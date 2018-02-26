<template lang="pug">
  div
    h2.mb-4 How to use this site
    p.
      Welcome to this <strong>Play by Post</strong> web app. It's used to play RPGs in a <i>post-by-post</i> format: each
      player posts what their character is doing as they have time during their day, and their story unfolds asynchronously,
      without all players having to meet physically or virtually at the same time to play.
    p.mt-4.
      This app is designed to facilitate that effort: user's posts are layed out in a simple unnested old-school forum layout,
      one above the other, with the most recent posts at the bottom.
    p.mt-4.
      As this board is intended to be viewed by anyone but only allow participating users to post, anyone can view the existing
      posts without having to have an account or log in, but only users with accounts can make additional posts to the board.
    h2 Layout
    p.mt-4 The layout of the site is as follows:
    ul
      li.
        At the top of the page is the <strong>toolbar</strong>. It contains, left to right, the sidenav button, the web app title,
        a blank area, and on the far right, an indicator telling you which character you are logged in as (if you're logged in).
      li.
        On the left side of the page is nothing special, until you click the <i>sidenav button</i> (also known as a <i>hamburger button</i>)
        on the toolbar. Once you do, the side navigation panel will slide out from the left of the page and show additional links:
      ul.ml-5
        li <strong>Posts</strong> - view the posts on the board (basically, go to the home page)
        li <strong>Help</strong> - view this page
        li <strong>Log in</strong> - log into your character/user account
        li <strong>Register</strong> - create a new account
        li <strong>Profile</strong> - change your name, password, email, etc.
        li <strong>Log out</strong> - log out of your account
        li <strong>Close this</strong> - close the side navigation panel
      span You can also close the panel by clicking or tapping anywhere on the shaded area that appears
      li.
        In the center of the screen lies that page's content, whether that's posts, the login form, this text, etc.
    h2.mt-4 A note on passwords
    p.mt-4.
      With all of the security concerns and leaks and such going on these days, I've made sure that I'm using industry best-practices
      for hashing and salting your password. I never see it. The database content for your password will be something unrecognizable
      and irreversible, like <strong>$2a$10$rBV2JDeWW3.vKyeQcM8fFO4777l4bVeQgDL6VIkxqlzQ7TCalQvla</strong>. If you have a software
      background or are just curious, this app is <i>open source</i> and can be viewed on
      <a href="https://github.com/Celeo/PlayByPost">my GitHub account</a>.
    h2.mt-4.mb-4 Posting
    p.
      New posts can be made by logged-in users by going to the <strong>Posts</strong> page and scrolling to the bottom. There, you will
      an area titled <strong>New post</strong>. The editor below that is a <i>WYSIWYG</i> editor: "What you see if what you get." It's
      more like a word processing program's editor than what you can see on some other play by post forums, like Paizo's.
    p.
      Type your post into the text area. You can use the buttons on the top of the editor to bold your text, add italics, add lists, and
      change the color to blue to denote "out of character" text. The rightmost button will remove all text formatting from whatever you
      have highlighted when you click it. Additionally, some common <i>keyboard shortcuts</i> work, like <strong>Ctrl+B</strong> to bold.
    h2.mt-4.mb-4 Rolling
    p.
      There is a widget titled "Dice Roller" that appears to the right of (or under, if on a small screen) the text widget where you write
      the content of your post. To roll some dice, first type in <i>what</i> you're rolling for, whether that's a Perception check, attacking,
      dealing damage, etc. Then, select the number of the <i>same type of dice</i> you're rolling, then the number of sides on that dice, and
      your net modifier, if any, to that roll.
    p.
      For a DnD 5e example, if you're rolling a Perception check, and are proficient in the skill and have a proficiency stat of 3 and Wisdom
      of 12, then your roll would be <code>1d20 + 3 + 1</code> which would be <code>1</code>, <code>20</code>, and <code>4</code> in the row.
      Then, click "Roll!". You'll see your dice and the result in the area under the button. If you have another check or attack or whatever
      to roll in the same post, repeat the process, substituting the new values for the next roll.
    p.
      If you have a roll that require multiple dice with multiple sides, then use the green plus button to the right of the first row to add
      additional rows. For another DnD 5e example, if you're rolling Investigation and have been given a bonus dice to roll alongside your d20,
      then you might set the first row in the roller to <code>1d20 + 3</code>, add a second row, and set that row to <code>1d4</code>. You
      can remove rows that you don't need with the red minus button to the left of the green button (only if you have more than one row).
    h2.mt-4.mb-4 Editing
    p.
      For 30 minutes after submiting your post, you can edit it. You can not edit your rolls, but you can change the text.
    p.
      When viewing posts, you'll see a small edit icon to the right of the submission date in the title of the post. Clicking on this icon
      will take you to a new page that allows you to change the content of the post and submit it. If you do so, the content of your post
      is changed for good. You can change your post as many times as you need to inside that editing window time, but after the window
      closes, your post is as it's going to be moving forward.
    h3.mt-2 Submitting
    p.
      Once you have written what your character does, click the blue <strong>Save</strong> button under the editor (it'll only be clickable
      if you've written something). You posting will be saved and shown on the board at the bottom.
    h2.mt-4.mb-4 Mobile
    p.
      This site should work on mobile. The site should recognize the viewing size of your device and shift the content and spacing accordingly.
    h2.mt-4.mb-4 Bugs / issues / requests
    p.
      You can submit issues and requests to <a href="mailto:celeodor@gmail.com">celeodor@gmail.com</a> or, if you use GitHub, the
      <a href="https://github.com/Celeo/PlayByPost">repo</a>.
</template>

<script>
export default {}
</script>

<style lang="stylus" scoped>
code
  background-color rgba(0, 0, 0, 0.05) !important
</style>
